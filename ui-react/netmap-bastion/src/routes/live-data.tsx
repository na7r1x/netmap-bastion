import BasePage from "../components/BasePage"
import { dateFormatAliases, EuiText, formatDate } from "@elastic/eui"
import Navigation from "../components/Navigation"

import { Chart, Settings, BarSeries, Axis, niceTimeFormatByDay, timeFormatter, ScaleType } from "@elastic/charts"
import { EUI_CHARTS_THEME_DARK, EUI_CHARTS_THEME_LIGHT } from '@elastic/eui/dist/eui_charts_theme';
import '@elastic/charts/dist/theme_only_light.css';
import '@elastic/charts/dist/theme_only_dark.css';

const headerProps = {
    pageTitle: "Live Data"
}

const isDarkTheme = false;
const euiTheme = isDarkTheme ? EUI_CHARTS_THEME_DARK.theme : EUI_CHARTS_THEME_LIGHT.theme;

const historyResponse = await fetch("http://localhost:8000/history");
const data = await historyResponse.json();

console.log(data)

const bottomPanel = (
    <Chart
        size={["100%", 500]}
    >
        <Settings
            theme={euiTheme}
            showLegend={false}
        />
        <BarSeries
            id="test"
            name="test"
            data={data}
            xScaleType={ScaleType.Time}
            yScaleType={ScaleType.Linear}
            xAccessor="time"
            yAccessors={["packetCount"]}

        />

        <Axis
            title="Time"
            // title={formatDate(Date.now(), dateFormatAliases.date)}
            id="bottom-axis"
            position="bottom"
            tickFormat={timeFormatter(niceTimeFormatByDay(1))}
            
            />
        <Axis
            title="Packet Count"
            id="left-axis"
            position="left"
            // showGridLines
            // tickFormat={(d) => Number(d).toFixed(2)}
        />
    </Chart>
)
export default function () {
    return (
        <BasePage sidebar={<Navigation />} header={headerProps} content={bottomPanel} />
    )
}