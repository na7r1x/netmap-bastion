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

const data = [
    { x: 0, y: 2 },
    { x: 1, y: 7 },
    { x: 2, y: 3 },
    { x: 3, y: 6 },
    { x: 4, y: 6 },
    { x: 5, y: 6 },
    { x: 6, y: 6 },
    { x: 7, y: 6 },
    { x: 8, y: 6 },
    { x: 9, y: 6 },
    { x: 10, y: 6 },
    { x: 11, y: 6 },
];

const bottomPanel = (
    <Chart
        size={{ height: 200 }}
    >
        <Settings
            theme={euiTheme}
            showLegend={false}
        />
        <BarSeries
            id="test"
            name="test"
            data={data}
            xScaleType={ScaleType.Linear}
            yScaleType={ScaleType.Linear}
            xAccessor={'x'}
            yAccessors={['y']}

        />

        <Axis
            title={formatDate(Date.now(), dateFormatAliases.date)}
            id="bottom-axis"
            position="bottom"
            tickFormat={timeFormatter(niceTimeFormatByDay(1))}

        />
        <Axis
            id="left-axis"
            position="left"
            showGridLines
            tickFormat={(d) => Number(d).toFixed(2)}
        />
    </Chart>
)
export default function () {
    return (
        <BasePage sidebar={<Navigation />} header={headerProps} content={bottomPanel} />
    )
}