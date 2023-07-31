# Event Management API

The Event Management API is a comprehensive solution for managing events, tickets, users, and host requests. It provides different roles with specific functionalities, including Guest, User, Host, Moderator, and Admin. The API allows users to perform various actions such as browsing and filtering events, purchasing tickets, requesting refunds, creating events, approving host requests, and more.

## Roles and Actions

### Guest Role

As a Guest, users have limited access to the system. They can perform the following actions:

- **✅ View Event (GET):** Retrieve detailed information about a specific event, including name, description, location, date, ticket availability, etc.

- **⏳ List Events (GET):** Retrieve a list of events with pagination and sorting options.
  - Filtering options:
    - Status: Ongoing, Upcoming, or Past events.
    - Duration: Events of a specific duration (e.g., 1-day, 2-day, etc.).
    - Month: Events occurring in a particular month.
    - City: Events taking place in a specific city.
    - Month and City: Events in a particular city during a specific month.
    - Minimum Tickets: Events with at least a certain number of tickets left.

- **⏳ Search Events (GET):** Search for events by name, description, or location.

### User Role

In addition to the actions available to Guests, Users have additional functionalities:

- **✅ Buy Tickets (POST):** Purchase tickets for an event. Users can buy multiple tickets at once, but they cannot buy tickets for ongoing events.
  
- **✅ Request to Become a Host (POST):** Users can request to become a host. If the request is denied, the user will not be able to request again for 30 days.

- **⏳ View Bought Tickets (GET):** Retrieve a list of all purchased tickets with pagination and sorting options.

- **⏳ Search Bought Tickets (GET):** Search for tickets by event name or location.

### Host Role

As Hosts, users have additional privileges on top of User actions:

- **✅ All User Role actions** are available to Hosts.

- **✅ Create Event (POST):** Create a new event with details such as name, description, location, date, etc.

- **✅ List Created Events (GET):** Retrieve a list of all events created by the host with pagination and sorting options.
  - ⏳ Filtering options:
    - Location and Date: Events in a specific location and date range.
    - Status: Ongoing or not ongoing events.
      
- **⏳ Update Event (PUT):** Update event information, including description (at any time) and name, location, date, etc. (only if no tickets have been sold).

- **⏳ Delete Event (DELETE):** Delete an event from the system, but only if no tickets have been sold.

- **⏳ Search Events (GET):** Search for events by name, description, location, date, etc.

- **⏳ View Tickets Sold for an Event (GET):** Retrieve a list of all tickets sold for a specific event with pagination and sorting options.

### Moderator Role

Moderators have the following actions available to them:

- **✅ View Requests to Become a Host (GET):** Retrieve a list of host requests with pagination and sorting options.

- **✅ Approve or Reject a Host Request (PUT):** Moderators can review and approve or reject host requests.

### Admin Role

Administrators have additional privileges and can perform the following actions:

- **✅ All Moderator Role actions** are available to Administrators.

- **⏳ Create Moderator (POST):** Create a new moderator account.

## Database Structure

The following tables are used in the database:

- **Users:** Stores user information including id, email, password, role, created_at, and password_updated_at.

- **Events:** Contains event details such as id, host_id (foreign key to users table), name, description, location, total_tickets, tickets_left, start_date, end_date, status, and created_at.

- **Tickets:** Stores ticket information including id, user_id, event_id, quantity, and created_at.

- **User_Host_Request:** Manages user requests to become a host with fields such as id, user_id (foreign key to users table), moderator_id (foreign key to users table), status, and created_at.

---

The Event Management API offers a robust set of features for managing events, tickets, user roles, and host requests. The provided roles allow for different levels of access and control within the system. By utilizing this API, you can easily build applications that empower users to browse, purchase, and manage event tickets, while also enabling hosts to create and manage events efficiently.
