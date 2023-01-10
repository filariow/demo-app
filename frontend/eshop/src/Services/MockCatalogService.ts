import { CatalogService, ProductModel, CatalogError } from './CatalogService'
import { v4 as uuidv4 } from 'uuid';

export class MockCatalogService implements CatalogService {
  products: ProductModel[]
  constructor () {
    this.products = []
    for (let index = 1; index <= 20; index++) {
      this.products.push(this.generateProduct())
    }
  }

  async GetProduct (id: string): Promise<CatalogError | ProductModel> {
    const pp = this.products.filter(x => x.id === id)
    if (pp.length > 0) return await Promise.resolve(pp[0])

    return await Promise.resolve({ Code: '404', Description: 'Product not found', Error: undefined })
  }

  async GetAllProducts (): Promise<CatalogError | ProductModel[]> {
    return await Promise.resolve(this.products)
  }

  private generateProduct (): ProductModel {
    const ru = Math.floor(Math.random() * 100)
    const id = uuidv4()
    
    return {
      id,
      name: `Mocked product - ${id}`,
      photoUrl: `https://picsum.photos/id/${ru}/1000`,
      unitSold: ru
    }
  }
}

export const DefaultMockCatalogService = new MockCatalogService()
