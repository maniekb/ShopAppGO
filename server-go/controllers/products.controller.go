package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"example/web-service-gin/services"
	"example/web-service-gin/models"
)

type ProductsController struct {
	productsService services.ProductsService
}

func NewProductsController(productsService services.ProductsService) ProductsController {
	return ProductsController{productsService}
}

func (pc *ProductsController) GetProducts(ctx *gin.Context) {
	var queryParams = ctx.Request.URL.Query()

	var searchTerm string = ""
	var manufacturers []int
	var priceFrom int
	var priceTo int

	if len(queryParams["searchTerm"]) != 0 {
		searchTerm = queryParams["searchTerm"][0]
	}
	if len(queryParams["manufacturer"]) != 0 {
		strs := strings.Split(queryParams["manufacturer"][0], ",")
		for i := range strs {
			var m, _ = strconv.Atoi(strs[i])
			manufacturers = append(manufacturers, m)
		}
	}
	if len(queryParams["priceFrom"]) != 0 {
		priceFrom, _ = strconv.Atoi(queryParams["priceFrom"][0])
	}
	if len(queryParams["priceTo"]) != 0 {
		priceTo, _ = strconv.Atoi(queryParams["priceTo"][0])
	}

	skip, err := strconv.ParseInt(queryParams["skip"][0], 10, 64)
	limit, err := strconv.ParseInt(queryParams["limit"][0], 10, 64)

	products, total, err := pc.productsService.GetProducts(searchTerm, skip, limit, manufacturers, priceFrom, priceTo)

	if err != nil {
		// if err == mongo.ErrNoDocuments {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "products": products, "total": total})
		// 	return
		// }
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if len(products) == 0 {
		products = make([]*models.ProductDBResponse, 0)
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "products": products, "total": total})
	return
}

func (pc *ProductsController) UploadProduct(ctx *gin.Context) {
	var product *models.ProductUploadInput

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	newProduct, err := pc.productsService.UploadProduct(product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "product": newProduct})
	return
}

func (pc *ProductsController) GetProduct(ctx *gin.Context) {
	var queryParams = ctx.Request.URL.Query()

	var id string

	if len(queryParams["id"]) != 0 {
		id = queryParams["id"][0]
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}
	

	product, err := pc.productsService.GetProduct(id)

	if err != nil {
		// if err == mongo.ErrNoDocuments {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "products": products, "total": total})
		// 	return
		// }
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "product": product})
	return
}




// func (pc *ProductsController) UploadImage(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("File Upload Endpoint Hit")

//     // Parse our multipart form, 10 << 20 specifies a maximum
//     // upload of 10 MB files.
//     r.ParseMultipartForm(10 << 20)
//     // FormFile returns the first file for the given key `myFile`
//     // it also returns the FileHeader so we can get the Filename,
//     // the Header and the size of the file
//     file, handler, err := r.FormFile("myFile")
//     if err != nil {
//         fmt.Println("Error Retrieving the File")
//         fmt.Println(err)
//         return
//     }
//     defer file.Close()
//     fmt.Printf("Uploaded File: %+v\n", handler.Filename)
//     fmt.Printf("File Size: %+v\n", handler.Size)
//     fmt.Printf("MIME Header: %+v\n", handler.Header)

//     // Create a temporary file within our temp-images directory that follows
//     // a particular naming pattern
//     tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
//     if err != nil {
//         fmt.Println(err)
//     }
//     defer tempFile.Close()

//     // read all of the contents of our uploaded file into a
//     // byte array
//     fileBytes, err := ioutil.ReadAll(file)
//     if err != nil {
//         fmt.Println(err)
//     }
//     // write this byte array to our temporary file
//     tempFile.Write(fileBytes)
//     // return that we have successfully uploaded our file!
//     fmt.Fprintf(w, "Successfully Uploaded File\n")
// }