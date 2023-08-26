import { createApp } from 'vue'
import { createPinia } from 'pinia'

import './style.css'
import App from './App.vue'

// Font Awesome Icons
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {fas} from "@fortawesome/free-solid-svg-icons";



/* add icons to the library */
library.add(fas)

const pinia = createPinia()
const app = createApp(App)

app
    .use(pinia)
    .component('font-awesome-icon', FontAwesomeIcon)
    .mount('#app')