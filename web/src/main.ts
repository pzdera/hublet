import '@fontsource-variable/inter';
import '@fontsource-variable/geist';
import '@fontsource-variable/manrope';
import '@fontsource-variable/ibm-plex-sans';

import { mount } from 'svelte';

import App from './App.svelte';

import './styles.css';
import './editor.css';
import './card-layout.css';
import './section-sortable.css';
import './resource-info.css';
import './service-resources.css';
import './appearance.css';
import './section-surfaces.css';
import './service-icons.css';
import './local-icon-manager.css';

mount(App, {
  target: document.getElementById('app')!
});
