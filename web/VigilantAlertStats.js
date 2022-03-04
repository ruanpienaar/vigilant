import React from 'react';
import rd3 from 'react-d3-library';
import node from './VigilantD3Chart';
const RD3Component = rd3.Component;

export default class VigilantAlertStats extends React.Component {
    state = {
        d3: ''
    }
    componentDidMount() {
        console.log('component did mount');
        //console.log(node);
        alert(node);
        this.setState( () => ( {d3: node} ));
    }
    render() {
        return (
            <div>
                d3
                <RD3Component data={this.state.d3} />
            </div>
        )
    }
};