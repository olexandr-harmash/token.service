<template>
  <div class="container">
    <div class="row">
      <div class="col-lg-8 col-md-8 col-sm-8 mx-auto">
        <svg class="bi" width="250" height="250" fill="blue">
          <use xlink:href="~/bootstrap-icons/bootstrap-icons.svg#emoji-smile" />
        </svg>
        <h2 class="mt-4"><strong>Key</strong></h2>
        <h3 style="color: #476582">Protected store for you'r data</h3>
        <div class="mt-4 log-buttons">
          <router-link class="login_button" to="/login">Log in</router-link>
        </div>
      </div>
    </div>
    <CardsGrid />
    <div class="row">
      <div class="col-lg-8 col-md-8 col-sm-8 mx-auto my-3 footer">
        <div class="row">
          <div class="col-lg-8 col-md-8 col-sm-8 mx-auto">
            <a href="#">1</a>
            <a href="#">2</a>
          </div>
        </div>
        <div class="row">
          <div class="col-lg-8 col-md-8 col-sm-8 mx-auto my-3">
            <a href="#" class="sign_in_button">To Sponsor</a>
          </div>
        </div>
      </div>
      <p>Nice Licensed | Copyright © 2019-present Key Chain St</p>
    </div>
  </div>
</template>

<script>
import CardsGrid from "/Users/iuada0h0/Desktop/auth_service/ui/src/components/CardsGrid.vue";

export default {
  components: {
    CardsGrid,
  },
};
</script>

<style lang="scss">
@import "/Users/iuada0h0/Desktop/auth_service/ui/src/assets/main.scss";

.footer {
  background-color: #f8f8f8;
  .row {
    display: flex;
    a {
      margin-left: 5%;
      margin-right: 5%;
      width: 10em;
    }
  }
}

.log-buttons {
  display: flex;
  flex-direction: row;
  justify-content: center;
}

.sign_in_button {
  @include btn-text;
  @include btn-sz;
  border: 1px solid black;
  &:hover {
    @include btn-hover;
  }
}

.login_button {
  @include btn-text;
  @include btn-sz;
  background-color: $btn_col;
  &:hover {
    @include btn-hover;
  }
}
</style>
