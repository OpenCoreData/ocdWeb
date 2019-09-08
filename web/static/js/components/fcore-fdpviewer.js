/*jshint esversion: 6 */

import {
    html,
    render
} from './lit-html.js';

(function () {
    class SimpleGet extends HTMLElement {
        constructor() {
            super();

            // Pull in the JSON GET boilerplate from the Fence code.
            // Nee to make sure to centralize the component work in geocomponents.org
            // How to do this locally?

            const resID = this.getAttribute('res-id');

            // CAUTION DEV / DEMO HACK..  comment out in production!!!!!!
            var newstr = resID.replace(/opencoredata.org/i, '192.168.2.89:9900');
            // console.log(newstr);

            // GET test
            function tj_providers(id) {
                return fetch(id, {
                    // headers: { 'Accept': 'application/ld+json', },
                })
                    .then(function (response) {
                        return response.json();
                    })
                    .then(function (myJson) {
                        console.log("=== fdp viewer ===")
                        console.log(myJson);
                        // console.log(JSON.stringify(myJson));
                        // return JSON.stringify(myJson);
                        return myJson;
                    });
            }

            // GET test call...
            tj_providers(newstr).then((feature) => {
                this.attachShadow({ mode: 'open' });

                var  count = Object.keys(feature.resources).length;
                const detailsTemplate = [];

                var i;
                for (i = 0; i < count; i++) {
                    detailsTemplate.push( html`<div><a target="_blank"
                    href="${feature.resources[i].path}">${feature.resources[i].name}</a></div>`);
                }

                var h = html`<div style="margin-top:10px">
                <span>${feature.title}</span><br>
                ${detailsTemplate}</div>`;
                // this.shadowRoot.innerHTML = `${h}`;
                render(h, this.shadowRoot);

            });
        }
    }
    window.customElements.define('fcore-fdpviewer', SimpleGet);
})();


