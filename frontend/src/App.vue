<script setup lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';


</script>

<template>
  <div class="container mx-auto max-w-screen-md">
    <!-- Logo -->
    <div class="w-full flex mt-24">
      <img src="./assets/logo-lightmode.svg" class="mx-auto" alt="TundraDNS Logo" />
    </div>
    
    <!-- Welcome Text -->
    <div class="w-full flex mt-16">
      <p class="text-primary">
        Welcome to the TundraDNS Portal! Our single-page web application empowers you to effortlessly oversee your active domains and their associated records. You currently have the following domain active, you can switch between your other domains via the drop down.
      </p>
    </div>

    <!-- Domain Box -->
    <div class="w-full flex mt-8">
      <div class="w-full h-24 flex bg-secondary rounded-lg justify-between">
        <div class="h-full flex ml-5">
          <div class="w-16 h-16 flex my-auto bg-primary rounded-lg">
            <font-awesome-icon icon="cloud" class="text-white my-auto mx-auto h-2/5" />
          </div>
          <div class="my-auto ml-8">
            <p class="text-xl  text-primary">
              tundra-dns.io
            </p>
            <div class="flex">
              <div class="border-2 border-success-full rounded-full bg-success-faint">
                <p class="text-success text-xs px-5 py-0.5">active</p>
              </div>
            </div>
          </div>
        </div>
        <div class="h-full flex">
          <font-awesome-icon icon="chevron-right" class="text-primary my-auto mr-10" size="xl"/>
        </div>
      </div>
    </div>

  </div>
</template>
