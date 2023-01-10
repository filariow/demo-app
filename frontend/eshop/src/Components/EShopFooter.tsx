import React from 'react'

function EShopFooter (): JSX.Element {
  return (
    <footer className="page-footer font-small d-flex flex-wrap justify-content-between align-items-center m-0 pt-3" style={{ background: '#1a1a20' }}>
      <div className="container">
        <div className="row">

          <div className="col-4 text-muted mb-0 justify-content-center pt-2">
            <p>Red Hat © 2022 Company, Inc</p>
          </div>

          <div className="offset-4 col-4">
            <ul className="nav flex-column justify-content-center pb-3" >
              <li className="nav-item">
                <a href="https://github.com/redhat-developer/service-binding-operator/" className="nav-link px-2 text-muted">Service Binding Operator</a>
              </li>
              <li className="nav-item">
                <a href="https://redhat-developer.github.io/service-binding-operator/userguide/intro.html" className="nav-link px-2 text-muted">Documentation</a>
              </li>
              <li className="nav-item">
                <a href="https://github.com/servicebinding/spec" className="nav-link px-2 text-muted">Specification</a>
              </li>
            </ul>
          </div>

        </div>
      </div>
    </footer >
  )
}

export default EShopFooter
