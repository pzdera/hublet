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
import './appearance.css';
import './service-icons.css';
import './local-icon-manager.css';

import './ui-polish.css';
import './section-surfaces.css';
import './dashboard-identity.css';
import './section-rhythm.css';
import './shortcuts.css';
import './weather.css';
import './command-search.css';
import './mobile-polish.css';

mount(App, {
  target: document.getElementById('app')!
});
