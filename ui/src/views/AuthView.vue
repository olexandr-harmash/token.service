<template>
  <div class="container">
    <div class="auth-wrapper">
      <form class="auth" onsubmit="return false">
        <div>
          <h1>Authorize</h1>
          <p>The client would like to perform actions on your behalf.</p>
        </div>

        <input v-model="code" type="code" placeholder="Code" required />
        <button
          @click="doAuth"
          type="submit"
          class="btn btn-primary btn-lg"
          style="width: 200px"
        >
          Allow
        </button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  code: "",
  methods: {
    doAuth() {
      if (this.emailLogin === "" || this.passwordLogin === "") {
        this.emptyFields = true;
      } else {
        alert("You are now logged in");
        fetch(`http://localhost:9096/auth?code=${this.code}`, {
          method: "GET",
          credentials: "include",
        })
          .then((responce) => {
            console.log(responce);
            alert(responce);
            if (responce.ok) {
              this.$router.push("/");
            }
          })
          .catch((err) => {
            alert(err);
          });
        //TODO redirect to client
      }
    },
  },
};
</script>

<style lang="scss">
.auth-wrapper {
  background-color: whitesmoke;
  margin-top: 0.5em;
}
.auth {
  height: 250px;
  text-align: left;
  display: flex;
  flex-direction: column;
  margin-left: 2.5em;
  * {
    margin-top: 0.5em;
  }
  input {
    width: 40%;
    margin-left: 25%;
  }
}
</style>
