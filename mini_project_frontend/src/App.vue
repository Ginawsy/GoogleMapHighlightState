<script setup>
import { ref } from 'vue'
import { loadScript } from "vue-plugin-load-script";

// Stores States' names queried from backend
const queried_states_with_prefix = ref([]);
// Stores a lambda for drawing/refreshing the map. This will be initialized
// after Google Map API script is loaded.
var drawRegionsMap;

// Loads Google Map API script and renders map
loadScript('https://www.gstatic.com/charts/loader.js')
  .then(() => {
    google.charts.load('current', {
        'packages':['geochart'],
      });
      drawRegionsMap = () => {
        var data_table = []
        if (queried_states_with_prefix.value.length == 1) {
          data_table = [[queried_states_with_prefix.value[0]]];
        }
        var data = google.visualization.arrayToDataTable([
          ['State'],
        ].concat(data_table));

        var options = {
          region: 'US',
          resolution: 'provinces'
        };

        var chart = new google.visualization.GeoChart(document.getElementById('regions_div'));

        chart.draw(data, options);
      }
      google.charts.setOnLoadCallback(drawRegionsMap);
  })
  .catch((e) => {
    console.log(e)
  })

// Queries backend for States starting with user's input inside the textbox.
// Rerenders the map after getting the State names.
// If there's only one State name returned, highlights the State on the map.
async function queryStatesWithPrefix(e) {
  queried_states_with_prefix.value = []
  fetch(`http://localhost:3000/query?state_prefix=${e.target.value}`)
    .then((res) => {
      return res.json();
    })
    .then((states_with_prefix) => {
      queried_states_with_prefix.value = states_with_prefix;
      drawRegionsMap();
    })
    .catch((err) => {
      console.log('Error: ', err)
    })
}

</script>

<template>
  <div id="regions_div" style="width: 900px; height: 500px;"></div>  
  <input list="States" placeholder="Please type in a State's name"
    style="width: 200px;" 
    @input="queryStatesWithPrefix"/>
  
  <datalist id="States">
      <option v-for="state in queried_states_with_prefix">{{state}}</option>
  </datalist>
</template>