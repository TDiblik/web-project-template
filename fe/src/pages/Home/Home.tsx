import {IfLoggedIn} from "../../components/IfLoggedIn";
import Layout from "../../components/Layout";

function Home() {
  return (
    <IfLoggedIn redirectToLogin={true}>
      <Layout>
        <div>
          <h1 className="text-3xl font-bold mb-6">Dashboard</h1>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div className="card bg-base-100 shadow-md p-6">
              <h2 className="text-xl font-semibold mb-2">Users</h2>
              <p className="text-3xl font-bold">1,245</p>
            </div>
            <div className="card bg-base-100 shadow-md p-6">
              <h2 className="text-xl font-semibold mb-2">Revenue</h2>
              <p className="text-3xl font-bold">$12,345</p>
            </div>
            <div className="card bg-base-100 shadow-md p-6">
              <h2 className="text-xl font-semibold mb-2">Orders</h2>
              <p className="text-3xl font-bold">245</p>
            </div>
          </div>

          <div className="mt-8">
            <h2 className="text-2xl font-bold mb-4">Recent Activity</h2>
            <ul className="space-y-2">
              <li className="p-4 bg-base-100 rounded-lg shadow-sm">New user registered: Jane Doe</li>
              <li className="p-4 bg-base-100 rounded-lg shadow-sm">Order #1234 completed</li>
              <li className="p-4 bg-base-100 rounded-lg shadow-sm">Revenue goal reached this month</li>
            </ul>
          </div>
        </div>
      </Layout>
    </IfLoggedIn>
  );
}

export default Home;
