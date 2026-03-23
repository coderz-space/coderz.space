import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Text, View, StyleSheet } from 'react-native';
import { Colors, Typography, Spacing } from '../theme';
import { MentorTabParamList } from '../types';

import MentorDashboardScreen from '../screens/mentor/DashboardScreen';
import DoubtsScreen from '../screens/mentor/DoubtsScreen';
import PollsScreen from '../screens/mentor/PollsScreen';
import MentorProfileScreen from '../screens/mentor/ProfileScreen';

const Tab = createBottomTabNavigator<MentorTabParamList>();

function TabIcon({ icon, focused, badge }: { icon: string; focused: boolean; badge?: number }) {
  return (
    <View style={[tabStyles.iconWrap, focused && tabStyles.iconWrapActive]}>
      <Text style={[tabStyles.icon, focused && tabStyles.iconActive]}>{icon}</Text>
      {badge && badge > 0 ? (
        <View style={tabStyles.badge}>
          <Text style={tabStyles.badgeText}>{badge > 9 ? '9+' : badge}</Text>
        </View>
      ) : null}
    </View>
  );
}

export default function MentorTabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={{
        headerShown: false,
        tabBarStyle: tabStyles.bar,
        tabBarShowLabel: true,
        tabBarLabelStyle: tabStyles.label,
        tabBarActiveTintColor: Colors.primary,
        tabBarInactiveTintColor: Colors.textDisabled,
      }}
    >
      <Tab.Screen
        name="DashboardTab"
        component={MentorDashboardScreen}
        options={{
          tabBarLabel: 'Mentees',
          tabBarIcon: ({ focused }) => <TabIcon icon="⊞" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="DoubtsTab"
        component={DoubtsScreen}
        options={{
          tabBarLabel: 'Doubts',
          tabBarIcon: ({ focused }) => <TabIcon icon="?" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="PollsTab"
        component={PollsScreen}
        options={{
          tabBarLabel: 'Polls',
          tabBarIcon: ({ focused }) => <TabIcon icon="◈" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="ProfileTab"
        component={MentorProfileScreen}
        options={{
          tabBarLabel: 'Profile',
          tabBarIcon: ({ focused }) => <TabIcon icon="◎" focused={focused} />,
        }}
      />
    </Tab.Navigator>
  );
}

const tabStyles = StyleSheet.create({
  bar: {
    backgroundColor: Colors.surface,
    borderTopWidth: 1,
    borderTopColor: Colors.surfaceBorder,
    height: 64,
    paddingBottom: Spacing.sm,
    paddingTop: Spacing.sm,
  },
  label: {
    ...Typography.labelSmall,
    fontSize: 10,
    marginTop: -2,
  },
  iconWrap: {
    width: 32,
    height: 32,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
    position: 'relative',
  },
  iconWrapActive: {
    backgroundColor: Colors.primaryMuted,
  },
  icon: {
    fontSize: 18,
    color: Colors.textDisabled,
  },
  iconActive: {
    color: Colors.primary,
  },
  badge: {
    position: 'absolute',
    top: -2,
    right: -2,
    backgroundColor: Colors.error,
    borderRadius: 8,
    width: 16,
    height: 16,
    alignItems: 'center',
    justifyContent: 'center',
  },
  badgeText: {
    ...Typography.caption,
    color: Colors.white,
    fontSize: 9,
    fontWeight: '800',
  },
});