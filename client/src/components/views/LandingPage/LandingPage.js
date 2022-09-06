import Axios from 'axios';
import { Card, Icon, Row, Col } from "antd";
import React, { useEffect, useState } from 'react'
import { PRODUCT_SERVER } from '../../Config.js';
import CheckBox from './Sections/CheckBox.js';
import RadioBox from './Sections/RadioBox';
import { price, manufacturers } from './Sections/Data'
import SearchFeature from './Sections/SearchFeature';
import { Link } from 'react-router-dom';

const { Meta } = Card;

function LandingPage() {

    const Limit = 8;
    const [Products, setProducts] = useState([])
    const [Skip, setSkip] = useState(0)
    const [ProductsTotal, setProductsTotal] = useState(0)
    const [SearchTerms, setSearchTerms] = useState("")
    const [Filters, setFilters] = useState({
        manufacturer: [],
        price: []
    })

    useEffect(() => {
        const variables = {
            skip: Skip,
            limit: Limit
        }
        getProducts(variables)
    }, [])

    const renderCars = Products.map((product, index) => {
            return <Col lg={6} md={8} xs={24}>
                <Link to={`/product/${product.id}`}>
                <Card 
                    hoverable={true} 
                >
                    <Meta title={product.title} description={`$${product.price}`}/>
                    
                </Card>
                </Link>
            </Col>
    })

    const getQueryParams = (variables) => {
        let queryParams = {
            skip: variables.skip,
            limit: variables.limit
        }

        if(variables.filters !== undefined){
            if(variables.filters.manufacturer !== undefined && variables.filters.manufacturer.length > 0) {
                queryParams["manufacturer"] = variables.filters.manufacturer.join(',')
            }
            if(variables.filters.price !== undefined && variables.filters.price.length > 0) {
                queryParams["priceFrom"] = variables.filters.price[0]
                queryParams["priceTo"] = variables.filters.price[1]
            }
            if(variables.searchTerm !== "") {
                queryParams["searchTerm"] = SearchTerms
            }
        }

        return queryParams
    }

    const getProducts = (variables) => {

        let queryParams = getQueryParams(variables)

        Axios.get(`${PRODUCT_SERVER}/getProducts`, { params: queryParams})
            .then(response => {
                if(response.data.success) {
                    if(variables.loadMore) {
                        setProducts([...Products, ...response.data.products])
                    } else {
                        setProducts(response.data.products)
                    }
                    setProductsTotal(response.data.total)
                } else {
                    alert('Failed to fetch products.')
                }
            })
    }

    const onLoadMore = () => {
        let skip = Skip + Limit;

        const variables = {
            skip: Skip,
            limit: Limit,
            loadMore: true
        }
        getProducts(variables)
        setSkip(skip)
    }

    const showFilteredResults = (filters) => {
        const variables = {
            skip: 0,
            limit: Limit,
            filters: filters
        }
        getProducts(variables)
        setSkip(0)
    }

    const handlePrice = (value) => {
        const data = price;
        let array = [];

        for (let key in data) {
            if(data[key]._id === parseInt(value, 10)) {
                array = data[key].array;
            }
        }
        return array
    }

    const handleFilters = (filters, category) => {
        const newFilters = { ...Filters }
        newFilters[category] = filters

        if(category === "price") {
            let priceValues = handlePrice(filters)
            newFilters[category] = priceValues
        }

        showFilteredResults(newFilters)
        setFilters(newFilters)
    }

    const updateSearchTerms = (newSearchTerm) => {

        const variables = {
            skip: 0,
            limit: Limit,
            filters: Filters,
            searchTerm: newSearchTerm
        }

        setSkip(0)
        setSearchTerms(newSearchTerm)

        getProducts(variables)
    }

    return (
        <div style={{ width: '75%', margin: '3rem auto'}}>
            <div style={{ textAlign: 'center' }}>
                <h2>Find your dream phone <Icon type="mobile" /></h2>
            </div>

            <Row gutter={[16, 16]}>
                <Col lg={12} xs={24} >
                    <CheckBox
                        list={manufacturers}
                        handleFilters={filters => handleFilters(filters, "manufacturer")}
                    />
                </Col>
                <Col lg={12} xs={24}>
                    <RadioBox
                        list={price}
                        handleFilters={filters => handleFilters(filters, "price")}
                    />
                </Col>
            </Row>

            <div style={{ display: 'flex', justifyContent: 'flex-end', margin: '1rem auto' }}>
                <SearchFeature
                    refreshFunction={updateSearchTerms}
                 />
            </div>

            {Products.length === 0 ? 
                <div style={{ display: 'flex', height: '300px', justifyContent: 'center', alignItems: 'center' }}>
                    <h2>No phones yet...</h2>
                </div> :
                <div>
                    <Row gutter={[16,16]}>
                        {renderCars}
                    </Row>

                </div>
            }
            <br/>
            <br/>

            {ProductsTotal >= Limit && 
                <div style={{ display: 'flex', justifyContent: 'center'}}>
                <button onClick={onLoadMore}>Load More</button>
                </div>
            }

            
        </div>
    )
}

export default LandingPage
