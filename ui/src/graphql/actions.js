import gql from 'graphql-tag'

/*
 *******************************************
 *************  User / Profile  ************
 *******************************************
 */
export const GET_CURRENT_USER = gql`
  query getCurrentUser {
    viewer {
      id
      firstName
      lastName
      email
      gravatar
      createdAt
      organizations {
        nodes {
          id
          name
        }
      }
      projects {
        nodes {
          id
          name
        }
      }
    }
  }
`

/*
 *******************************************
 *************  Organizations  *************
 *******************************************
 */

export const CREATE_ORGANIZATION = gql`
  mutation createOrganization($organization: Organization!) {
    createOrganization(organization: $organization) {
      id
      name
      description
    }
  }
`

export const UPDATE_ORGANIZATION = gql`
  mutation updateOrganization($organization: Organization!) {
    updateOrganization(organization: $organization) {
      id
      name
      description
    }
  }
`

export const DELETE_ORGANIZATION = gql`
  mutation deleteOrganization($id: ID!) {
    deleteOrganization(id: $id) {
      id
    }
  }
`

export const ADD_MEMBER_TO_ORGANIZATION = gql`
  mutation addMemberToOrganization($organization: Organization!, $member: OrganizationMember!) {
    addMemberToOrganization(organization: $organization, member: $member) {
      organization {
        id
        name
        members
      }
    }
  }
`

export const REMOVE_MEMBER_FROM_ORGANIZATION = gql`
  mutation removeMemberFromOrganization(
    $organization: Organization!
    $member: OrganizationMember!
  ) {
    removeMemberFromOrganization(organization: $organization, member: $member) {
      organization {
        id
        name
        members
      }
    }
  }
`

/*
 *******************************************
 ***************  Projects  ****************
 *******************************************
 */

export const GET_PROJECTS = gql`
  query {
    viewer {
      projects {
        nodes {
          id
          name
        }
      }
    }
  }
`

/*
 *******************************************
 **************  Local state  **************
 *******************************************
 */

export const UPDATE_SOME_STATE = gql`
  mutation updateSomeState($value: String) {
    updateSomeState(value: $value) @client
  }
`
