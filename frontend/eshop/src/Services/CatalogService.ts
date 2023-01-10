export interface CatalogService {
  GetProduct: (id: string) => Promise<ProductModel | CatalogError>
  GetAllProducts: () => Promise<ProductModel[] | CatalogError>
}

export interface ProductModel {
  id: string
  name: string
  photoUrl: string
  unitSold: number
}

export class CatalogError {
  Code: string
  Description: string
  Error: CatalogError | undefined

  constructor (code: string, description: string, error: CatalogError | undefined) {
    this.Code = code
    this.Description = description
    this.Error = error
  }
}
