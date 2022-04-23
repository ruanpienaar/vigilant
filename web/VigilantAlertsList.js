import React from 'react';
import VigilantAlert from './VigilantAlert';
const VigilantAlertsList = (props) => (
    <ul>
    {
        props.alerts.map((alert) => (
            <VigilantAlert key={alert.generatorURL} alert={alert} />
        ))
    }
    </ul>
);
export default VigilantAlertsList;