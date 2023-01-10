import 'bootstrap/dist/css/bootstrap.css' // Import precompiled Bootstrap css
import '@fortawesome/fontawesome-free/css/all.css'
import './ProductDetail.css'
import { ProductModel, CatalogError } from '../Services/CatalogService'
import { useNavigate, useParams } from 'react-router-dom'
import { Button, Form, InputGroup, Modal } from 'react-bootstrap'
import React, { Dispatch, SetStateAction, useEffect, useMemo, useState } from 'react'
import { DefaultHTTPCatalogService } from '../Services/HTTPCatalogService'
import { DefaultHTTPOrdersService } from '../Services/HTTPOrdersService'
import { OrdersError } from '../Services/OrdersService'

interface ProductDetailState {
  id: string | undefined

  // router
  navigate: Function

  // API request
  response: CatalogError | ProductModel | undefined
  setResponse: Dispatch<SetStateAction<CatalogError | ProductModel | undefined>>
  error: any
  setError: Dispatch<SetStateAction<any>>

  // dialog
  show: boolean
  setShow: Dispatch<SetStateAction<boolean>>
  showError: boolean
  setShowError: Dispatch<SetStateAction<boolean>>
  quantity: string
  setQuantityValue: Dispatch<SetStateAction<string>>
  quantityError: string
  setQuantityError: Dispatch<SetStateAction<string>>
}

function ProductDetail (): JSX.Element {
  const { id } = useParams()
  const [response, setResponse] = useState<CatalogError | ProductModel | undefined>(undefined)
  const [error, setError] = useState<any>(undefined)
  const [show, setShow] = useState(false)
  const [showError, setShowError] = useState(false)
  const [quantity, setQuantityValue] = useState<string>('1')
  const [quantityError, setQuantityError] = useState('')
  const navigate = useNavigate()

  const state: ProductDetailState = useMemo(() => {
    return {
      id,
      navigate,
      response,
      setResponse,
      error,
      setError,
      show,
      setShow,
      showError,
      setShowError,
      quantity,
      setQuantityValue,
      quantityError,
      setQuantityError
    }
  }, [id, error, navigate, quantity, quantityError, response, show, showError])

  useEffect(() => {
    if (id === undefined) {
      state.setError('Error: no id provided')
      return
    }

    DefaultHTTPCatalogService.GetProduct(id)
      .then(v => state.setResponse(v))
      .catch(err => state.setError(err))
  }, [id])

  const rp = function (): JSX.Element {
    if (state.error !== undefined) {
      return Error(state.error)
    }

    if (state.response != null) {
      return ProductSection(state)
    }

    return Error('Bug: neither error, nor response')
  }

  return (
    <div className="ProductDetail">
      <div className="ProductDetail-title-background py-4">
        <h1 className="ProductDetail-title">Product Detail</h1>
      </div>
      <div className="row pt-4 px-4 m-0">
        {rp()}
      </div>
    </div >
  )
}

function ProductSection (state: ProductDetailState): JSX.Element {
  const resp = state.response
  if (resp instanceof CatalogError) {
    const e = resp
    return (<div><p>Error {e.Code}: {e.Description}</p></div>)
  }

  const p = resp as unknown as ProductModel
  return (<div className="offset-4 col-4 py-4 px-5" key={p.id}>{ProductCard(state, p)}</div>)
}

function ProductCard (state: ProductDetailState, p: ProductModel): JSX.Element {
  return (
    <div className="card" style={{ color: 'black' }}>
      <img className="card-img-top" src={p.photoUrl} alt="Card cap" />
      <div className="card-body">
        <h3 className="card-title">{p.name}</h3>
        <h5 className="card-text">Units Sold: {p.unitSold}</h5>
        {BuyButtonWithModal(state, p)}
      </div>
    </div>
  )
}

function BuyButtonWithModal (state: ProductDetailState, p: ProductModel): JSX.Element {
  const handleClose = (): void => state.setShow(false)
  const handleShow = (): void => state.setShow(true)

  const handleCloseError = (): void => state.setShowError(false)
  const handleShowError = (): void => state.setShowError(true)

  let error: string | undefined

  const setQuantity = (v: string): void => {
    state.setQuantityError('')
    state.setQuantityValue(v)
  }

  const getQuantity = (): number | undefined => {
    const i = parseInt(state.quantity)
    return !isNaN(i) ? i : undefined
  }

  const buy = (): void => {
    const qi = getQuantity()
    if (qi === undefined) {
      state.setQuantityError('Quantity must be a positive integer')
      return
    }

    const req = DefaultHTTPOrdersService.PlaceOrder([
      {
        id: p.id,
        name: p.name,
        photoUrl: p.photoUrl,
        unitsOrdered: qi
      }
    ])

    req.then(res => {
      handleClose()

      if (res instanceof OrdersError) {
        error = res.Description
        handleShowError()
        return
      }

      state.navigate(`/orders/${res.id}`, { replace: false })
    }).catch(err => {
      error = err
      handleShowError()
    })
  }

  return (
    <div className="pt-3">
      <Button variant="primary" onClick={handleShow}>
        Buy Product
      </Button>

      {/* Buy modal */}
      <Modal show={state.show} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Buy product &lsquo;{p.name}&rsquo;</Modal.Title>
        </Modal.Header>

        <Modal.Body>
          <p>Do you really want to buy the product &lsquo;{p.name}&rsquo;</p>

          <InputGroup className="mb-3">
            <InputGroup.Text >📦</InputGroup.Text>
            <Form.Control
              placeholder="Quantity"
              aria-label="Quantity"
              value={state.quantity}
              onChange={e => setQuantity(e.target.value)}
              type="number"
            />
          </InputGroup>
          <Form.Text muted>{state.quantityError}</Form.Text>

        </Modal.Body>

        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>Cancel</Button>
          <Button variant="primary" onClick={buy}>Buy</Button>
        </Modal.Footer>
      </Modal>

      {/* Error modal */}
      <Modal show={state.showError} onHide={handleCloseError}>
        <Modal.Header closeButton>
          <Modal.Title>Buy Product {p.name}</Modal.Title>
        </Modal.Header>

        <Modal.Body>
          <p>Error while placing order: {error}</p>
        </Modal.Body>

        <Modal.Footer>
          <Button variant="secondary" onClick={handleCloseError}>Ok</Button>
        </Modal.Footer>
      </Modal>
    </div>
  )
}

function Error (message: string): JSX.Element {
  return (<div><p>{message}</p></div>)
}

export default ProductDetail
