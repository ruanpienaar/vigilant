import React from 'react';
import 'react-date-range/dist/styles.css'; // main css file
import 'react-date-range/dist/theme/default.css'; // theme css file
import { DateRangePicker } from 'react-date-range';
export const VigilantDatePicker = (props) => {
    return (
        <div>
            <DateRangePicker
                ranges={[props.selectionRange]}
                onChange={props.handleSelect}
            />
            <button onClick={props.showState}>showState</button>
        </div>
    );
}
export default VigilantDatePicker;