import { mount } from 'svelte';

import App from './App.svelte';

import './styles.css';
import './editor.css';
import './card-layout.css';
import './section-sortable.css';

mount(App, {
  target: document.getElementById('app')!
});
