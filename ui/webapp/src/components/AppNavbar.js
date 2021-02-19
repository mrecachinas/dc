import { Link } from "react-router-dom";
import { Dropdown, Nav, Navbar } from "rsuite";

export default function AppNavbar({ message, updateTime }) {
  return (
    <Navbar>
      <Navbar.Header>
        <a
          href="#"
          className="navbar-brand logo"
          style={{ padding: "18px 20px", display: "inline-block" }}
        >
          DC
        </a>
      </Navbar.Header>
      <Navbar.Body>
        <Nav>
          <Dropdown title="Status" trigger={["hover", "click"]} componentClass={Link} to="/status">
            <Dropdown.Item componentClass={Link} to="/active">
              Active Status
            </Dropdown.Item>
            <Dropdown.Item componentClass={Link} to="/historical">
              Historical Status
            </Dropdown.Item>
          </Dropdown>
          <Nav.Item componentClass={Link} to="/status">
            Status
          </Nav.Item>
          <Nav.Item componentClass={Link} to="/tasks">
            Tasks
          </Nav.Item>
        </Nav>
        <Nav pullRight>
          {message ? 
          <Nav.Item disabled={true}>
            {message}
          </Nav.Item> : null
          }
          {/* <Nav.Item icon={<Icon icon="cog" />} >Settings</Nav.Item> */}
        </Nav>
      </Navbar.Body>
    </Navbar>
  );
}
