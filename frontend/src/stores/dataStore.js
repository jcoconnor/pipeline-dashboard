import { observable } from 'mobx'
import Train from './train';

const axios = require('axios');
const moment = require('moment');

class Product {
  name: ""
  pipeline: ""
  wallClockTime: 0
  totalTimeDuration: 0
  startTime: ""
  errors: 0
  transients: 0
  allJobs: []

  constructor(product, jobs) {
    this.name = product.Name;
    this.pipeline = product.Pipeline;
    this.startTime = product.StartTime;
    this.wallClockTime = product.WallClockTime;
    this.totalTimeDuration = product.TotalTimeDuration;
    this.errors = product.Errors;
    this.transients = product.Transients;
    this.allJobs = jobs;
  }

  GetPipelines() {
    var retVal = [];
    retVal = this.allJobs.filter((job) => {
      return (this.pipeline === job.pipeline);
    });

    return retVal;
  }
}


class Job {
  url: ""
  pipeline: ""
  wallClockTime: 0
  totalTimeDuration: 0
  allTrains: []
  pipelineJob: ""
  startTime: ""
  endTime: ""
  errors: 0
  transients: 0
  version: ""
  buildNumber: 0

  constructor(job, trains) {
    this.url = job.URL;
    this.pipeline = job.Pipeline;
    this.pipelineJob = job.PipelineJob;
    this.wallClockTime = job.WallClockTime;
    this.totalTimeDuration = job.TotalTimeDuration;
    this.version = job.Version;
    this.jobDataStrings = job.JobDataStrings;
    this.buildNumber = job.BuildNumber;
    this.startTime = moment(job.JobDataStrings.StartTime, "YYYY-MM-DD HH:mm:ss Z PDT");
    this.endTime = moment(job.JobDataStrings.EndTime, "YYYY-MM-DD HH:mm:ss Z PDT");
    this.errors = job.Errors;
    this.transients = job.Transients;
    this.allTrains = trains;
  }

  totalFormatted() {
    return `${this.jobDataStrings.TotalHours}H, ${this.jobDataStrings.TotalMinutes}M`
  }

  wallClockFormatted() {
    return `${this.jobDataStrings.WallClockTimeHours}H, ${this.jobDataStrings.WallClockTimeMinutes}M`
  }

  GetTrains() {
    var retVal = [];
    retVal = this.allTrains.filter((train) => {
      return ((this.pipeline === train.pipeline) && (this.version === train.version));
    });

    return retVal;

  }
}

class Link {
  title: ""
  url: ""

  constructor(product, jobs) {
    this.url = product.URL;
    this.title = product.Title;
  }
}

class DataStore {
  @observable data = 'supercalifragilisticexpialidocious'
  @observable products = []
  @observable jobs = []
  @observable trains = []
  @observable links = []
  @observable title = ""
  @observable state = ""

  fetchProducts(cb) {
    var store = this;
    axios.get('/api/1/products')
      .then((res: any) => res.data)
      .then(function(res: any) {
        store.trains   = res.Trains.map((train)     => new Train(train))
        store.jobs     = res.Jobs.map((job)         => new Job(job, store.trains))
        store.products = res.Products.map((product) => new Product(product, store.jobs))
        store.links    = res.Links.map((link) => new Link(link))
        store.title    = res.Title
        store.state    = "done"
        cb()
      })
      .catch((err: any) => {
        console.log("in axios ", err)
        store.state = "error"
      })

  }
}

export default DataStore
