import React from 'react';
import VigilantDatePicker from './VigilantDatePicker';
import VigilantAlertStats from './VigilantAlertStats';
import VigilantAlertsList from './VigilantAlertsList';
import Axios from 'axios';
export default class VigilantApp extends React.Component {
    state = {
        startDate: new Date(),
        endDate: new Date(),
        alerts: []
    }
    showState = () => {
        console.log(this.state);
    }
    handleDateSelect = (s) => {
        // console.log(s); // native Date object
        console.log( 'startDate '+s.selection.startDate );
        console.log( 'endDate '+s.selection.endDate );
        const alerts = this.httpGetAlerts(s.selection.startDate, s.selection.endDate);
        if (alerts != undefined || alerts != null) {
            this.setState( () => ({
                alerts,
                startDate: s.selection.startDate,
                endDate: s.selection.endDate
            }) );
        }
    }
    componentDidMount(){
        // TODO: how would we get the appropriate HOST?
        // .window methods might give you a public facing url, which is no good...
        // for now hard-code (localhost), to allow local-dev to carry on.
        const alerts = this.httpGetAlerts(this.state.startDate, this.state.endDate);
        if (alerts != undefined || alerts != null) {
            this.setState( () => ({
                alerts
            }) );
        }
    }
    render(){
        const selectionRange = {
            startDate: this.state.startDate,
            endDate: this.state.endDate,
            key: 'selection',
        };
        //alert('start: '+this.state.startDate + ' -> end: ' +this.state.endDate);
        try {
            console.log('# Alerts fetched '+this.state.alerts.length+' for start: '+this.state.startDate + ' -> end: ' +this.state.endDate);
        } catch  (e) {
            console.log("")
        }
        return (
            <div>
                <h1>Vigilant</h1>
                <button>Choose a daterage</button>
                <div>
                    <VigilantDatePicker selectionRange={selectionRange} handleDateSelect={this.handleDateSelect} showState={this.showState} />
                    <VigilantAlertsList alerts={this.state.alerts} />
                </div>
            </div>
        );
    }

    httpGetAlerts(startDate, endDate){
        Axios.get(
            'http://localhost:8801/api/list/all-alerts',
            { params: {
                startDate: startDate,
                endDate: endDate
            }}
        )
        .then(
            (response) =>
                {
                    console.log({
                        action: 'check-alerts',
                        alerts: response.Alerts
                    });
                    if (response.Alerts != undefined || response.Alerts != null) {
                        return response.Alerts;
                    } else {
                        return [];
                    }
               }
        );
    }

}

//Axious.
 // const alerts = [1, 2];


// <VigilantAlertStats />