import React, { useState } from 'react';
import { EuiIcon, EuiSideNav, slugify } from '@elastic/eui';
import { useNavigate, useLocation } from 'react-router-dom'

import { appendIconComponentCache } from '@elastic/eui/es/components/icon/icon';

import { icon as EuiIconGauge } from '@elastic/eui/es/components/icon/assets/vis_gauge';
import { icon as EuiIconTimeline } from '@elastic/eui/es/components/icon/assets/timeline';
import { icon as EuiIconWrench } from '@elastic/eui/es/components/icon/assets/wrench';

// One or more icons are passed in as an object of iconKey (string): IconComponent
appendIconComponentCache({
    visGauge: EuiIconGauge,
    timeline: EuiIconTimeline,
    wrench: EuiIconWrench,
});

export default () => {
    const navigate = useNavigate();
    const location = useLocation();

    const [isSideNavOpenOnMobile, setIsSideNavOpenOnMobile] = useState(false);
    const [selectedItemName, setSelectedItem] = useState(location.pathname);

    const toggleOpenOnMobile = () => {
        setIsSideNavOpenOnMobile(!isSideNavOpenOnMobile);
    };

    const selectItem = (id: string) => {
        navigate(id)
        setSelectedItem(id);
    };

    const createItem = (name: string, data = {}) => {
        // NOTE: Duplicate `name` values will cause `id` collisions.
        const id = '/' + slugify(name)
        return {
            id: id,
            name,
            isSelected: selectedItemName === id,
            onClick: () => selectItem(id),
            ...data,
        };
    };

    const sideNav = [
        createItem('Monitoring', {
            onClick: undefined,
            icon: <EuiIcon type="visGauge" />,
            items: [
                createItem('Live Data'),
                createItem('Historic Data'),
                createItem('Alerting')
            ],
        }),
        createItem('Management', {
            onClick: undefined,
            icon: <EuiIcon type="timeline" />,
            items: [
                createItem('Agents'),
                createItem('Relays'),
            ],
        }),
        createItem('Settings', {
            onClick: undefined,
            icon: <EuiIcon type="wrench" />,
            items: [
                createItem('Account settings'),
            ],
        }),
    ];

    return (
        <EuiSideNav
            aria-label="Navigation"
            mobileTitle="Navigation"
            toggleOpenOnMobile={toggleOpenOnMobile}
            isOpenOnMobile={isSideNavOpenOnMobile}
            items={sideNav}
            style={{ width: 192 }}
        />
    );
};