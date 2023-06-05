import React from 'react'
import 'bootstrap/dist/css/bootstrap.css' // Import precompiled Bootstrap css
import '@fortawesome/fontawesome-free/css/all.css'
import './Home.css'
import { Link } from 'react-router-dom'

function Home (): JSX.Element {
  return (
    <div className="Home">
      <header className="Home-header" style={{ marginTop: '-104px', marginBottom: '-104px' }}>
        <p>
          Welcome to <span className="eshop-title">eShop</span>
        </p>
        <p style={{ paddingTop: '2rem' }}>
          As Warehouse Manager of the <i>Sapiens Inc</i><span className="ms-2">🦍</span>
          <br />
          <span>you have access to our <Link className="app-link" to={'catalog'}>Catalog</Link> to buy new products</span>
        </p>
        <p>
          <span>Visit the <Link className="app-link" to={'orders'}>Orders</Link> page to review already placed orders</span>
        </p>
      </header>
    </div>
  )
}

export default Home
