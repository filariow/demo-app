import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import { LinkContainer } from 'react-router-bootstrap'
import { Link } from 'react-router-dom'
import React from 'react'

function EShopNavbar (): JSX.Element {
  return (
    <div>
      <Navbar collapseOnSelect expand="lg" bg="dark" variant="dark">
        <Container>
          <Navbar.Brand as={Link} to="/" className='eshop-title'>📦 eShop</Navbar.Brand>
          <Navbar.Toggle aria-controls="responsive-navbar-nav" />
          <Navbar.Collapse id="responsive-navbar-nav">
            <Nav className="me-auto">
              <LinkContainer to="catalog">
                <Nav.Link>Catalog</Nav.Link>
              </LinkContainer>
              <LinkContainer to="orders">
                <Nav.Link>Orders</Nav.Link>
              </LinkContainer>
            </Nav>
            <Nav className="me-end">
              <Nav.Link href="https://github.com/redhat-developer/service-binding-operator/">
                <span className="pe-2">Sapiens Inc<span className="ms-2">🦍</span></span>
              </Nav.Link>
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>
    </div>
  )
}

export default EShopNavbar
