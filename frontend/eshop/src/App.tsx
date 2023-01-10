import React from 'react'
import { Routes, Route, Outlet } from 'react-router-dom'
import './App.css'
import EShopFooter from './Components/EShopFooter'
import EShopNavbar from './Components/EShopNavbar'
import Home from './Pages/Home'
import Catalog from './Pages/Catalog'
import Orders from './Pages/Orders'
import OrderDetail from './Pages/OrderDetail'
import ProductDetail from './Pages/ProductDetail'

function App (): JSX.Element {
  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="catalog" element={<Catalog />} />
          <Route path="catalog/:id" element={<ProductDetail />} />
          <Route path="orders" element={<Orders />} />
          <Route path="orders/:id" element={<OrderDetail />} />
        </Route>
      </Routes>
    </div>
  )
}

function Layout (): JSX.Element {
  return (
    <div>
      <header>
        <EShopNavbar />
      </header>
      <article>
        <Outlet />
      </article>
      <footer>
        <EShopFooter />
      </footer>
    </div>
  )
}

export default App
