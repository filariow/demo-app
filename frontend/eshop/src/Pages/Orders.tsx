import 'bootstrap/dist/css/bootstrap.css' // Import precompiled Bootstrap css
import '@fortawesome/fontawesome-free/css/all.css'
import './Orders.css'
import { OrderedProductModel, OrderModel, OrdersError } from '../Services/OrdersService'
import { useNavigate } from 'react-router-dom'
import { DefaultHTTPOrdersService } from '../Services/HTTPOrdersService'
import React, { useEffect, useState } from 'react'
import { Table } from 'react-bootstrap'

function Orders (): JSX.Element {
  const [response, setResponse] = useState<OrdersError | OrderModel[] | undefined>(undefined)
  const [error, setError] = useState<any>(undefined)
  const navigate = useNavigate()

  useEffect(() => {
    DefaultHTTPOrdersService.GetAllOrders()
      .then(v => setResponse(v))
      .catch(err => { setError(err) })
  }, [])

  const rp = (): JSX.Element => {
    if (error !== undefined) {
      return (<div>Error: {error}</div>)
    }

    if (response instanceof OrdersError) {
      const e = response
      return (<div><p>Error {e.Code}: {e.Description}</p></div>)
    }

    if (response != null) {
      // return OrdersSection(response, navigate)
      return (
        <div className="row pt-4 mt-1">
          <div className="col-8 offset-2">
            {OrdersTable(response, navigate)}
          </div>
        </div>
      )
    }

    return (<div>Error: Bug: neither error, nor response</div>)
  }

  return (
    <div className="Orders">
      <div className="Orders-title-background py-4">
        <h1 className="Orders-title">Orders</h1>
      </div>
      <div className="row m-0">
        {rp()}
      </div>
    </div>
  )
}

function OrdersTable (orders: OrderModel[], navigate: Function): JSX.Element {
  return (
    <Table striped bordered hover variant="dark">
      <thead>
        <tr>
          <th>Id</th>
          <th>Date</th>
          <th>Product</th>
          <th>Quantity</th>
        </tr>
      </thead>
      <tbody>
        {OrderTableBodyContent(orders, navigate)}
      </tbody>
    </Table>
  )
}

function OrderTableBodyContent (orders: OrderModel[], navigate: Function): JSX.Element[] {
  const p = (o: OrderModel): OrderedProductModel | undefined =>
    o.orderedProducts.length === 0 ? undefined : o.orderedProducts[0]

  return orders.map((o: OrderModel) =>
    <tr onClick={() => { navigate(`/orders/${o.id}`, { replace: false }) }}
      key={o.id} style={{ cursor: 'pointer' }}>
      <td>{o.id}</td>
      <td>{new Date(o.date).toLocaleDateString()} {new Date(o.date).toLocaleTimeString()}</td>
      <td>{p(o)?.name ?? 'No product'}</td>
      <td>{p(o)?.unitsOrdered ?? ''}</td>
    </tr >)
}

export default Orders
