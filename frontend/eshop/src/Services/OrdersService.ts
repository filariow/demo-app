export interface OrdersService {
  GetOrder: (id: string) => Promise<OrderModel | OrdersError>
  GetAllOrders: () => Promise<OrderModel[] | OrdersError>
  PlaceOrder: (pp: OrderedProductModel[]) => Promise<OrderModel | OrdersError>
}

export interface OrderModel {
  id: string
  date: Date
  orderedProducts: OrderedProductModel[]
}

export interface OrderedProductModel {
  id: string
  name: string
  photoUrl: string
  unitsOrdered: number
}

export class OrdersError {
  Code: string
  Description: string
  Error: Error | undefined

  constructor (code: string, description: string, error: Error | undefined) {
    this.Code = code
    this.Description = description
    this.Error = error
  }
}
