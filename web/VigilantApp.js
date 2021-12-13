import React from 'react';
import VigilantDatePicker from './VigilantDatePicker';
export default class VigilantApp extends React.Component {
    state = {
        startDate: new Date(),
        endDate: new Date()
    }
    showState = () => {
        console.log(this.state);
    }
    handleSelect = (s) => {
        // console.log(s); // native Date object
        console.log( s.selection.startDate );
        console.log( s.selection.endDate );
        this.setState( () => (
            {
                startDate: s.selection.startDate,
                endDate: s.selection.endDate
            }
        ) )
    }
    render(){
        const selectionRange = {
            startDate: this.state.startDate,
            endDate: this.state.endDate,
            key: 'selection',
        };
        return (
            <div>
                <h1>Vigilant</h1>
                <button>Choose a daterage</button>
                <div>
                    <VigilantDatePicker selectionRange={selectionRange} handleSelect={this.handleSelect} showState={this.showState} />
                </div>
            </div>
        );
    }
}