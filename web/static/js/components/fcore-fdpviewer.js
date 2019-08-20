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

                var h = html`<div>${detailsTemplate}</div>`;
                // this.shadowRoot.innerHTML = `${h}`;
                render(detailsTemplate, this.shadowRoot);
                // this.shadowRoot.innerHTML = html `<div style="overflow-wrap: break-word;width=100%">
                //    Description: ${feature.description} <br>
                //    ${detailsTemplate}
                // </div> `;

                // var inputs = feature["resources"];
                // for (var i = 0; i < inputs.length; i++) {
                //     this.shadowRoot.innerHTML = `<div style="overflow-wrap: break-word;width=100%">
                //     File: ${inputs[i].name} <br>
                //     </div>`;
                // }

            });
        }
    }
    window.customElements.define('fcore-fdpviewer', SimpleGet);
})();


