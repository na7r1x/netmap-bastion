import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import "@elastic/eui/dist/eui_theme_light.css";
import "@elastic/charts/dist/theme_light.css";

import { EuiProvider, EuiThemeProvider } from '@elastic/eui';

import { RouteProvider } from './routes/RouteProvider'
import Navigation from './components/Navigation';

const router = createBrowserRouter(RouteProvider.getRoutes())

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <EuiProvider>
      <EuiThemeProvider colorMode='light'>
        <RouterProvider router={router} />
      </EuiThemeProvider>
    </EuiProvider>
  </React.StrictMode>,
)
