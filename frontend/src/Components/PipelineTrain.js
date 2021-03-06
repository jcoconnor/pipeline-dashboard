import React, { Component } from 'react'
import { observer } from 'mobx-react'

import Moment from 'react-moment';

@observer
class PipelineTrain extends Component {

  hms(nanoseconds) {
    var startSeconds = nanoseconds / 1000000000
    var hours = Math.floor(startSeconds / 3600);
    startSeconds = startSeconds - hours * 3600;

    var minutes = Math.floor(startSeconds / 60);
    // var seconds = startSeconds - minutes * 60;

    return `${hours}H, ${minutes}M`
  }

  openPipeline() {

  }

  render () {
    return (
      <div className="row train-row text-left">
        <div className="col-12">
          <a className="train-name" href={this.props.train.url} target="_blank" rel="noopener noreferrer">{this.props.train.name}</a>
        </div>
        <div className="col-5">
        </div>
        <div className="col-2">
          {this.props.train.queueTimeFormatted()}
        </div>
        <div className="col-2">
          {this.props.train.wallClockTime}
        </div>
        <div className="col-2">
          {this.props.train.durationFormatted()}
        </div>
        <div className="col-1">
          <Moment parse="YYYY-MM-DD HH:mm:ss Z PDT" format="YYYY/MM/DD HH:mm">{this.props.train.startTime}</Moment>
        </div>
        <div className="col-1">
          <Moment parse="YYYY-MM-DD HH:mm:ss Z PDT" format="YYYY/MM/DD HH:mm">{this.props.train.endTime}</Moment>
        </div>
        <div className="col-1">
          {this.props.train.errors} / {this.props.train.transients}
        </div>
      </div>
    )
  }
}

export default PipelineTrain
