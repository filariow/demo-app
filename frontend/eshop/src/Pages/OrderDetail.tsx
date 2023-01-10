import 'bootstrap/dist/css/bootstrap.css' // Import precompiled Bootstrap css
import '@fortawesome/fontawesome-free/css/all.css'
import './OrderDetail.css'
import { OrderModel, OrderedProductModel, OrdersError } from '../Services/OrdersService'
import { useParams } from 'react-router-dom'
import React, { useEffect, useState } from 'react'
import { DefaultHTTPOrdersService } from '../Services/HTTPOrdersService'
import { Card, ListGroup, ListGroupItem } from 'react-bootstrap'

function OrderDetail (): JSX.Element {
  return (
    <div className="OrderDetail" >
      {OrderTitle()}
      <header className="OrderDetail-header">
        <div className="container m-0">
          <div className="row">
            {OrderSection()}
          </div>
        </div>
      </header>
    </div>
  )
}

function OrderSection (): JSX.Element {
  const [response, setResponse] = useState<OrdersError | OrderModel | undefined>(undefined)
  const [error, setError] = useState<any>(undefined)
  const { id } = useParams()

  useEffect(() => {
    if (id === undefined) {
      setError('Error: no id provided')
      return
    }

    DefaultHTTPOrdersService.GetOrder(id)
      .then(v => setResponse(v))
      .catch(err => { setError(err) })
  }, [id])

  const rp = (): JSX.Element => {
    if (error !== undefined) {
      return Error(error)
    }

    if (response instanceof OrdersError) {
      const e = response
      return Error(e.Description)
    }

    if (response != null) {
      return OrderCard(response)
    }

    return Error('Bug: neither error, nor response')
  }

  if (response instanceof OrdersError) {
    const e = response
    return Error(`Error ${e.Code}: ${e.Description}`)
  }

  return (<div className="col-8 offset-2 py-4 px-5">{rp()}</div>)
}

function OrderTitle (): JSX.Element {
  return (
    <div className="Catalog-title-background py-4">
      <h1 className="Catalog-title">Order</h1>
    </div>
  )
}

function OrderCard (p: OrderModel): JSX.Element {
  return (
    <div>
      <Card>
        <Card.Body>
          <Card.Title>Order {p.id}</Card.Title>
          <Card.Text>
            <span style={{ fontSize: '1rem' }}>
              Placed on {new Date(p.date).toLocaleDateString()} at {new Date(p.date).toLocaleTimeString()}
            </span>
          </Card.Text>
        </Card.Body>
        <ListGroup className="ProductListing">
          {ProductsListing(p.orderedProducts)}
        </ListGroup>
      </Card>
    </div>
  )
}

function ProductsListing (pp: OrderedProductModel[]): JSX.Element[] {
  return pp.map(p => (
    <ListGroupItem style={{ backgroundColor: '#2c3034', color: 'white' }} key={p.id}>
      {ProductDetail(p)}
      </ListGroupItem>
  ))
}

function ProductDetail (p: OrderedProductModel): JSX.Element {
  return (
    <div>
      <div className="row">
        <div className="col-6">
          <span style={{ fontSize: '1rem', fontWeight: 'bold' }}>Name</span>
        </div>
        <div className="col-6">
          <span style={{ fontSize: '1rem' }}>{p.name}</span>
        </div>
      </div>
      <div className="row">
        <div className="col-6">
          <span style={{ fontSize: '1rem', fontWeight: 'bold' }}>Quantity</span>
        </div>
        <div className="col-6">
          <span style={{ fontSize: '1rem' }}>{p.unitsOrdered}</span>
        </div>
      </div>
    </div>
  )
}

function Error (message: string): JSX.Element {
  return (<div><p>{message}</p></div>)
}

export default OrderDetail
