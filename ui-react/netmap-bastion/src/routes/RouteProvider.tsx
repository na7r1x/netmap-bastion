import LiveData from "./live-data"
import Agents from "./agents"
import { Navigate } from 'react-router-dom'

class RouteProvider {
    static routes: Array<object> = [
        {
            path: "/",
            element: <Navigate to='/live-data' replace />
        },
        {
            path: "/live-data",
            element: <LiveData />,
        },
        {
            path: "/agents",
            element: <Agents />,
        },
    ];

    static getRoutes() {
        return this.routes;
    }
}

export { RouteProvider }