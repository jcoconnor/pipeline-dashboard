import { observable } from 'mobx'
import Train from './train';
import Product from './Product';
import Job from './Job';

const axios = require('axios');

class Link {
  title = ""
  url = ""

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
  @observable lastupdated = "-"

  fetchProducts() {
    var store = this;
    return new Promise((resolve, reject) => {
      axios({
        method: 'get',
        url: '/api/1/products',
        responseType: 'json'
      })
      .then((res) => {
        store.trains      = res.data.Trains.map((train)     => new Train(train))
        store.jobs        = res.data.Jobs.map((job)         => new Job(job, store.trains))
        store.products    = res.data.Products.map((product) => new Product(product, store.jobs))
        store.links       = res.data.Links.map((link)       => new Link(link))
        store.title       = res.data.Title
        store.lastupdated = res.data.LastUpdated
        store.state       = "done"
        console.log("###################RES################")
        console.log(res.data)
        resolve(store)
      })
      .catch((err) => {
        console.log("in axios ", err)
        store.state = "error"
        reject(err)
      })
    })
  }
}

export default DataStore
