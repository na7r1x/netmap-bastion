import BasePage from "../components/BasePage"
import { EuiText } from "@elastic/eui"
import Navigation from "../components/Navigation"

let headerProps = {
    pageTitle: "Agents"
}
export default function () {
    return (
        <BasePage sidebar={<Navigation />} header={headerProps} content={<EuiText>Hello World!</EuiText>} />
    )
}