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

        Axios.get(
            'http://localhost:8801/api/list/all-alerts',
            { params: {
                startDate: s.selection.startDate,
                endDate: s.selection.endDate
            }}
        )
        .then(
            (response) =>
                {
                    console.log(response);
                    //console.log(response.data);
                    //const jsonData = JSON.parse(response);
                    //alerts = response.data.Alerts;
                    // console.log('alerts:'+alerts);
                    if (response.Alerts != undefined || response.Alerts != null) {
                        this.setState( () => ( {
                            alerts: response.Alerts,
                            startDate: s.selection.startDate,
                            endDate: s.selection.endDate
                        } ));
                    }
               }
        );
        // this.setState( () => (
        //     {
        //         startDate: s.selection.startDate,
        //         endDate: s.selection.endDate
        //     }
        // ) );
    }
    componentDidMount(){
        // why were we doing this again ?
        //this.setState( () => ( {alerts: []} ));

        // TODO: how would we get the appropriate URL?
        // .window methods might give you a public facing url, which is no good...
        // for now hard-code, to allow local-dev to carry on.
        // TODO: Add date params to URL below.

        // TODO: read date from state, and add to URL below

        Axios.get('http://localhost:8801/api/list/all-alerts').then(
            (response) =>
                {
                    console.log(response.Alerts);
                    //console.log(response.data);
                    //const jsonData = JSON.parse(response);
                    //alerts = response.Alerts;
                    //console.log('alerts:'+alerts);
                    let alertsResponse;
                    if (response.Alerts != undefined || response.Alerts != null) {
                        alertsResponse = response.Alerts;
                        this.setState( () => ( {alerts: alertsResponse} ));
                    }
               }
        );
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
}

//Axious.
 // const alerts = [1, 2];


// <VigilantAlertStats />