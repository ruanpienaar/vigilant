import React from 'react';

import 'react-date-range/dist/styles.css'; // main style file
import 'react-date-range/dist/theme/default.css'; // theme css file
import { DateRangePicker } from 'react-date-range';

export default class VigilantApp extends React.Component {
    state = {
        timeFrame: '24h'
    }
    handleSelect(date){
        console.log(date); // native Date object
    }
    render(){
        const selectionRange = {
            startDate: new Date(),
            endDate: new Date(),
            key: 'selection',
        }
        return (
            <div>
                <p>Choose a daterage</p>
                <div>
                <DateRangePicker
                    ranges={[selectionRange]}
                    onChange={this.handleSelect}
                />
                </div>
            </div>
        );
    }
}