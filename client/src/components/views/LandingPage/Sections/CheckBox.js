import React, { useState } from "react";
import { Checkbox, Collapse } from "antd";

const {Panel} = Collapse

const manufacturers = [
    { "_id": 1, "name": "Mazda" },
    { "_id": 2, "name": "Audi" },
    { "_id": 3, "name": "Mercedes" },
    { "_id": 4, "name": "Volvo" },
    { "_id": 5, "name": "BMW" },
    { "_id": 6, "name": "Saab" }
]

function CheckBox(props) {

    const [Checked, setChecked] = useState([])

    const handleToggle = (value) => {
        const currentIndex = Checked.indexOf(value);
        const newChecked = [...Checked];

        if(currentIndex === -1) {
            newChecked.push(value)
        }
        else {
            newChecked.splice(currentIndex, 1)
        }

        setChecked(newChecked)
        props.handleFilters(newChecked)
    }

    const renderCheckboxLists = () => manufacturers.map((value, index) => (
        <React.Fragment key={index}>
            <Checkbox
                onChange={() => handleToggle(value._id)}
                type="checkbox"
                checked={Checked.indexOf(value._id) === -1 ? false : true }
            />
            <span>{value.name}</span>
        </React.Fragment>
    ))

    return (
        <div>
            <Collapse defaultActiveKey={['0']}>
                <Panel header="Manufacturers" key="1">
                    {renderCheckboxLists()}
                </Panel>
            </Collapse>
        </div>
    )
}

export default CheckBox