import React from 'react';
const VigilantAlert = (props) => (
    <li key={props.alert.GeneratorURL} >{props.alert.GeneratorURL}</li>
);
export default VigilantAlert;