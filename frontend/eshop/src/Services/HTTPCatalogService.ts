import { CatalogService, ProductModel, CatalogError } from './CatalogService'
import { Product, ProductApi, ProductApiFactory } from '../Clients/catalog'
import { AxiosResponse } from 'axios'

export class HTTPCatalogService implements CatalogService {
  apiClient: ProductApi

  constructor () {
    const url = process.env.REACT_APP_ESHOP_CATALOG_URL
    this.apiClient = ProductApiFactory(undefined, url) as ProductApi
  }

  async GetAllProducts (): Promise<CatalogError | ProductModel[]> {
    const resp: AxiosResponse<Product[], CatalogError> = await this.apiClient.getProducts()
    if (resp.status < 200 || resp.status > 299) {
      const error: CatalogError = new CatalogError(
        resp.status.toString(), resp.data.toString(), undefined)
      return error
    }

    return resp.data as ProductModel[]
  }

  async GetProduct (id: string): Promise<ProductModel | CatalogError> {
    const resp: AxiosResponse<Product | CatalogError> = await this.apiClient.getProductById(id)
    if (resp.status < 200 || resp.status > 299) {
      return resp.data as CatalogError
    }

    return resp.data as ProductModel
  }
}

export const DefaultHTTPCatalogService = new HTTPCatalogService()
