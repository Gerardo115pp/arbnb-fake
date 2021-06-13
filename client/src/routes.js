import Reservations from './pages/Reservations/Reservation.svelte';
import ManagerHome from './pages/ManagerHome/manager_home.svelte';
import UserProfile from './pages/UserProfile/UserProfile.svelte';
import SignUp from './pages/SignUp/SignUp.svelte';
import Login from './pages/login/Login.svelte';



const routes = {
    "/": Login,
    "/sign-up": SignUp,
    "/home-manager/:sk": ManagerHome,
    "/reservations/:sk": Reservations,
    "/user-profile/:sk": UserProfile
}

export { routes }