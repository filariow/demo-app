import { CatalogError, ProductModel } from './CatalogService'
import { DefaultMockCatalogService, MockCatalogService } from './MockCatalogService'
import { OrdersService, OrderModel, OrdersError, OrderedProductModel } from './OrdersService'
import { v4 as uuidv4 } from 'uuid';

export class MockOrdersService implements OrdersService {
  orders: OrderModel[] = []

  async PlaceOrder (pp: OrderedProductModel[]): Promise<OrderModel | OrdersError> {
    const id = uuidv4()
    const o: OrderModel = {
      id: id,
      date: new Date(),
      orderedProducts: pp
    }

    this.orders.push(o)

    for (const p of pp) {
      const cp = DefaultMockCatalogService.products.find(cp => cp?.id === p.id)
      if (cp !== undefined) { cp.unitSold += p.unitsOrdered }
    }

    return await new Promise(() => o)
  }

  private async init (): Promise<void> {
    for (let index = 1; index <= 8; index++) {
      this.orders[index] = await this.generateOrder()
    }
  }

  async GetOrder (id: string): Promise<OrdersError | OrderModel> {
    if (this.orders.length === 0) { await this.init() }

    const oo = this.orders.filter(x => x.id === id)
    if (oo.length > 0) return await Promise.resolve(oo[0])

    return await new Promise(() => { return { Code: '404', Description: 'Order not found', Error: undefined } })
  }

  async GetAllOrders (): Promise<OrdersError | OrderModel[]> {
    if (this.orders.length === 0) { await this.init() }
    return await new Promise(() => this.orders)
  }

  private async generateOrder (): Promise<OrderModel> {
    const catalog: MockCatalogService = DefaultMockCatalogService
    const pp = await catalog.GetAllProducts()
    const id = uuidv4()

    if (pp instanceof CatalogError) {
      return { id, date: new Date(), orderedProducts: [] }
    }

    // const np = Math.floor(Math.random() * (pp.length - 2)) + 1
    const np = 1

    return {
      id,
      date: new Date(),
      orderedProducts: this.getRandomProducts(np, pp)
    }
  }

  private getRandomProducts (n: number, pp: ProductModel[]): OrderedProductModel[] {
    const indices: number[] = []
    for (let index = 0; index < n; index++) {
      const ni = this.getRandomUnpickedNumber(pp.length - 1, indices)
      indices.push(ni)
    }

    const orderedProducts: OrderedProductModel[] = []
    for (const i of indices) {
      const p = pp[i]
      const q = Math.floor(Math.random() * 99) + 1
      orderedProducts.push({
        id: p.id,
        name: p.name,
        photoUrl: p.photoUrl,
        unitsOrdered: q
      })
    }

    return orderedProducts
  }

  private getRandomUnpickedNumber (max: number, picked: number[]): number {
    while (true) {
      const ni = Math.floor(Math.random() * max)
      if (picked.findIndex(x => x === ni) === -1) { return ni }
    }
  }
}

export const DefaultMockOrdersService = new MockOrdersService()
