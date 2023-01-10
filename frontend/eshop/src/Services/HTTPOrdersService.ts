import { AxiosResponse } from 'axios'
import { Order, OrderApi, OrderApiFactory, OrderedProduct } from '../Clients/orders'
import { OrdersService, OrderModel, OrdersError } from './OrdersService'

export class HTTPOrdersService implements OrdersService {
  apiClient: OrderApi

  constructor () {
    const url = process.env.REACT_APP_ESHOP_ORDERS_URL
    this.apiClient = OrderApiFactory(undefined, url) as OrderApi
  }

  async PlaceOrder (pp: OrderedProduct[]): Promise<OrderModel | OrdersError> {
    const resp: AxiosResponse<Order> = await this.apiClient.createOrder({
      orderedProducts: pp.map(p => {
        return {
          id: p.id,
          name: p.name,
          photoUrl: p.photoUrl,
          unitsOrdered: p.unitsOrdered
        }
      })
    })

    if (resp.status < 200 || resp.status > 299) {
      return resp.data as OrdersError
    }

    return resp.data as OrderModel
  }

  async GetOrder (id: string): Promise<OrdersError | OrderModel> {
    const resp: AxiosResponse<Order | OrdersError, any> = await this.apiClient.getOrderById(id)
    if (resp.status < 200 || resp.status > 299) {
      return resp.data as OrdersError
    }

    return resp.data as OrderModel
  }

  async GetAllOrders (): Promise<OrdersError | OrderModel[]> {
    const resp: AxiosResponse<Order[] | OrdersError> = await this.apiClient.getOrders()
    if (resp.status < 200 || resp.status > 299) {
      return resp.data as OrdersError
    }

    return resp.data as OrderModel[]
  }
}

export const DefaultHTTPOrdersService = new HTTPOrdersService()
