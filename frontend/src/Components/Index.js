import React, { Component } from 'react'
import { inject, observer } from 'mobx-react'

import Product from './Product';

import { OverlayTrigger } from 'react-bootstrap';
import { Icon } from '@puppet/react-components';

const renderWallClockTimeTooltip = props => (
  <div
    {...props}
    style={{
      backgroundColor: 'rgba(0, 0, 0, 0.85)',
      padding: '2px 10px',
      color: 'white',
      borderRadius: 3,
      ...props.style,
    }}
  >
    The time the first job starts and the last job stops.  Wall Clock time is not aggregated for products because there's too much variability
  </div>
);

const renderTotalTimeTooltip = props => (
  <div
    {...props}
    style={{
      backgroundColor: 'rgba(0, 0, 0, 0.85)',
      padding: '2px 10px',
      color: 'white',
      borderRadius: 3,
      ...props.style,
    }}
  >
    The time it would take to run all jobs consecutively.  This is meant to be a measure of how much "work" your job is doing.
  </div>
);



@inject('rootStore')
@observer
class Index extends Component {
  static isPrivate = true

  render() {
    var productsTable = [];

    if (this.props.rootStore.dataStore.products) {
      this.props.rootStore.dataStore.products.forEach((product, index) => {
        productsTable.push(
          <Product product={product} key={index} trains={this.props.rootStore.dataStore.trains} jobs={this.props.rootStore.dataStore.jobs} />
        )
      });
    }

    return (
      <div>
        <br />
        <h1>CI Dashboard</h1>
        <div className="row rc-table-header ">
          <div className="rc-table-header-cell col-3">
            <b>Name</b>
          </div>
          <div className="rc-table-header-cell col-1">
            <b>Build Number</b>
          </div>
          <div className="rc-table-header-cell col-1">
            <b>Version</b>
          </div>
          <div className="rc-table-header-cell col-1">
            <OverlayTrigger
              placement="bottom"
              delay={{ show: 250, hide: 400 }}
              overlay={renderWallClockTimeTooltip}
            >
              <span><b>Queue Time</b><Icon type="info-circle"></Icon></span>
            </OverlayTrigger>
          </div>
          <div className="rc-table-header-cell col-1">
            <OverlayTrigger
              placement="bottom"
              delay={{ show: 250, hide: 400 }}
              overlay={renderWallClockTimeTooltip}
            >
              <span><b>Wall Clock Time</b><Icon type="info-circle"></Icon></span>
            </OverlayTrigger>
          </div>
          <div className="rc-table-header-cell col-1">
            <OverlayTrigger
              placement="bottom"
              delay={{ show: 250, hide: 400 }}
              overlay={renderTotalTimeTooltip}
            >
              <span><b>Total Time</b><Icon type="info-circle"></Icon></span>
            </OverlayTrigger>
          </div>
          <div className="rc-table-header-cell col-1">
            <b>Start Time</b>
          </div>
          <div className="rc-table-header-cell col-1">
            <b>End Time</b>
          </div>
          <div className="rc-table-header-cell col-1">
            <b>Errors / Transients</b>
          </div>
          <div className="rc-table-header-cell col-1">
            <b>Detail</b>
          </div>
        </div>
        <div className="rc-table">
          {productsTable}
        </div>
        <ul>
          <li>Last Updated: {this.props.rootStore.dataStore.lastupdated}</li>
        </ul>
        <ul>
          {this.props.rootStore.dataStore.links.map((link) => {
            return <li><a href={link.url}>{link.title}</a></li>
          })}
        </ul>
      </div>
    )
  }
}

export default Index
