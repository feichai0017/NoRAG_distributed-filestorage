import React from 'react'
import ReactDOM from 'react-dom/client'
import store from './store'
import './index.scss'
import {Provider} from "react-redux";
import {RouterProvider} from "react-router-dom";
import router from "@/router/index.jsx";

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <Provider store={store}>
        <RouterProvider router={router}/>
    </Provider>
  </React.StrictMode>,
)
