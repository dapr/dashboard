export interface MenuItem {
  title: string;
  icon: string;
  link: string;
  name: string;
}

export const OVERVIEW_MENU_ITEM = {
  title: 'Overview',
  icon: 'home',
  link: '/overview',
  name: 'overview',
};

export const COMPONENTS_MENU_ITEM = {
  title: 'Components',
  icon: 'apps',
  link: '/components',
  name: 'components',
};

export const CONFIGURATIONS_MENU_ITEM = {
  title: 'Configurations',
  icon: 'build',
  link: '/configurations',
  name: 'configurations',
};

export const CONTROLPLANE_MENU_ITEM = {
  title: 'Control Plane',
  icon: 'settings',
  link: '/controlplane',
  name: 'status',
};

export let MENU_ITEMS: MenuItem[] = [
  OVERVIEW_MENU_ITEM,
];
