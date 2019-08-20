/*jshint esversion: 6 */

import {
    html,
    render
} from './lit-html.js';

 
(function () {
    class SimpleGet extends HTMLElement {
        constructor() {
            super();

            var obj;
            var inputs = document.getElementsByTagName('script');
            for (var i = 0; i < inputs.length; i++) {
                if (inputs[i].type.toLowerCase() == 'application/ld+json') {
                    obj = JSON.parse(inputs[i].innerHTML);
                }
            }

            // console.log(inputs);
            // console.log(obj);

            // var  count = Object.keys(obj).length;
            const detailsTemplate = [];
           
            detailsTemplate.push( html `<table>`);
            var item, subitem;
            for (item in obj) {
                // console.log( item + "   " + obj[item]);  
                detailsTemplate.push( html `
                <tr>
                     <td>${item}</td><td> ${obj[item]} </td>
                </tr>
                `)  

                // for (subitem in obj[item]) {
                //     console.log(obj[item][subitem]);      
                // }                
            }
            detailsTemplate.push( html `</table>`);



            var h = html`
                <div style="overflow-wrap: break-word;width=100%">
                    Feature: <a href="${obj["@id"]}"> 
                    ${obj["http://opencoredata.org/voc/csdco/v1/hole_ID"]}</a><br>
                    PI(s): ${obj["http://opencoredata.org/voc/csdco/v1/pi"]}<br><br>
                     ${obj["http://opencoredata.org/voc/csdco/v1/country"]} > 
                     ${obj["http://opencoredata.org/voc/csdco/v1/county_Region"]} > 
                     ${obj["http://opencoredata.org/voc/csdco/v1/location"]}
                     <br>
                     <a target="_blank" href="https://www.google.com/maps/search/?api=1&zoom=4&basemap=terrain&query=${obj["http://www.w3.org/2003/01/geo/wgs84_pos#lat"]},${obj["http://www.w3.org/2003/01/geo/wgs84_pos#long"]}">
                     (lat:  ${obj["http://www.w3.org/2003/01/geo/wgs84_pos#lat"]}
                      long:  ${obj["http://www.w3.org/2003/01/geo/wgs84_pos#long"]}</a>
                    )

                    <hr>

                      ${detailsTemplate}

                    </a>
                </div> `;

            this.attachShadow({ mode: 'open' });
            render(h, this.shadowRoot);
        }

    }
    window.customElements.define('geodex-displayfeature', SimpleGet);
})();

