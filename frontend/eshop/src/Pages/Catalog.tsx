import 'bootstrap/dist/css/bootstrap.css' // Import precompiled Bootstrap css
import '@fortawesome/fontawesome-free/css/all.css'
import './Catalog.css'
import { ProductModel, CatalogError } from '../Services/CatalogService'
import { useNavigate } from 'react-router-dom'
import { DefaultHTTPCatalogService } from '../Services/HTTPCatalogService'
import React, { useEffect, useState } from 'react'

function Catalog (): JSX.Element {
  const [response, setResponse] = useState<CatalogError | ProductModel[] | undefined>(undefined)
  const [error, setError] = useState<any>(undefined)
  const navigate = useNavigate()

  useEffect(() => {
    DefaultHTTPCatalogService.GetAllProducts()
      .then(v => setResponse(v))
      .catch(err => { setError(err) })
  }, [])

  const rp = (): JSX.Element => {
    if (error !== undefined) {
      return <div>Error: {error}</div>
    }

    if (response != null) {
      return ProductsSection(response, navigate)
    }

    return (<div>Error: Bug: neither error, nor response</div>)
  }

  return (
    <div className="Catalog">
      <div className="Catalog-title-background py-4">
        <h1 className="Catalog-title">Catalog</h1>
      </div>
      <div className="row pt-4 px-4 m-0 pb-4">
        {rp()}
      </div>
    </div >
  )
}

function ProductsSection (response: CatalogError | ProductModel[], navigate: Function): JSX.Element {
  if (response instanceof CatalogError) {
    const e = response
    return (<div><p>Error {e.Code}: {e.Description}</p></div>)
  }

  return (<>{ProductsCards(response, navigate)}</>)
}

function ProductsCards (response: ProductModel[], navigate: Function): JSX.Element[] {
  return response.map((p: ProductModel) =>
    <div className="col-3 py-4 px-5" key={p.id}>{ProductCard(p, navigate)}</div>)
}

function ProductCard (p: ProductModel, navigate: Function): JSX.Element {
  return (
    <div className="card" style={{ color: 'black', cursor: 'pointer' }}
      onClick={() => navigate(`/catalog/${p.id}`, { replace: false })}>

      <img className="card-img-top" src={p.photoUrl} alt="Card cap" />
      <div className="card-body">
        <h3 className="card-title">{p.name}</h3>
        <h5 className="card-text">Units Sold: {p.unitSold}</h5>
      </div>
    </div>
  )
}

export default Catalog
