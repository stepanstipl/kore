/**
 * Kore API
 * Kore API provides the frontend API (kore.appvia.io)
 *
 * The version of the OpenAPI document: 0.0.1
 * Contact: info@appvia.io
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';

/**
 * The TypesWhoAmI model module.
 * @module model/TypesWhoAmI
 * @version 0.0.1
 */
class TypesWhoAmI {
    /**
     * Constructs a new <code>TypesWhoAmI</code>.
     * @alias module:model/TypesWhoAmI
     */
    constructor() { 
        
        TypesWhoAmI.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>TypesWhoAmI</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/TypesWhoAmI} obj Optional instance to populate.
     * @return {module:model/TypesWhoAmI} The populated <code>TypesWhoAmI</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new TypesWhoAmI();

            if (data.hasOwnProperty('email')) {
                obj['email'] = ApiClient.convertToType(data['email'], 'String');
            }
            if (data.hasOwnProperty('teams')) {
                obj['teams'] = ApiClient.convertToType(data['teams'], ['String']);
            }
            if (data.hasOwnProperty('username')) {
                obj['username'] = ApiClient.convertToType(data['username'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getEmail() {
        return this.email;
    }

    /**
     * @param {String} email
     */
    setEmail(email) {
        this['email'] = email;
    }
/**
     * @return {Array.<String>}
     */
    getTeams() {
        return this.teams;
    }

    /**
     * @param {Array.<String>} teams
     */
    setTeams(teams) {
        this['teams'] = teams;
    }
/**
     * @return {String}
     */
    getUsername() {
        return this.username;
    }

    /**
     * @param {String} username
     */
    setUsername(username) {
        this['username'] = username;
    }

}

/**
 * @member {String} email
 */
TypesWhoAmI.prototype['email'] = undefined;

/**
 * @member {Array.<String>} teams
 */
TypesWhoAmI.prototype['teams'] = undefined;

/**
 * @member {String} username
 */
TypesWhoAmI.prototype['username'] = undefined;






export default TypesWhoAmI;

