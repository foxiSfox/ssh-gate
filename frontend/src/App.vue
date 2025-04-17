<template>
  <div>
    <nav class="nav">
      <div class="nav-container">
        <h1 class="nav-title">SSH Gate</h1>
        <div class="nav-links">
          <RouterLink class="nav-link" activeClass="active" to="/users">Users</RouterLink>
          <RouterLink class="nav-link" activeClass="active" to="/servers">Servers</RouterLink>
          <RouterLink class="nav-link" activeClass="active" to="/access">Access</RouterLink>
        </div>
      </div>
    </nav>

    <main class="container">
      <router-view />
    </main>
  </div>
</template>
