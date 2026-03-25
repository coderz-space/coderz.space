import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Text, View, StyleSheet } from 'react-native';
import { Colors, Typography, Spacing } from '../theme';
import { MenteeTabParamList } from '../types';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

import MenteeDashboardScreen from '../screens/mentee/DashboardScreen';
import LeaderboardScreen from '../screens/mentee/LeaderboardScreen';
import ProfileScreen from '../screens/mentee/ProfileScreen';

const Tab = createBottomTabNavigator<MenteeTabParamList>();

function TabIcon({ icon, focused }: { icon: string; focused: boolean }) {
  return (
    <View style={[tabStyles.iconWrap, focused && tabStyles.iconWrapActive]}>
      <Text style={[tabStyles.icon, focused && tabStyles.iconActive]}>
        {icon}
      </Text>
    </View>
  );
}

export default function MenteeTabNavigator() {
  const insets = useSafeAreaInsets();
  return (
    <Tab.Navigator
      screenOptions={{
        headerShown: false,
        tabBarStyle: {
          ...tabStyles.bar,
          height: 64 + insets.bottom, // ADD insets.bottom
          paddingBottom: insets.bottom + 4, // ADD insets.bottom
        },
        tabBarShowLabel: true,
        tabBarLabelStyle: tabStyles.label,
        tabBarActiveTintColor: Colors.primary,
        tabBarInactiveTintColor: Colors.textDisabled,
      }}
    >
      <Tab.Screen
        name="DashboardTab"
        component={MenteeDashboardScreen}
        options={{
          tabBarLabel: 'Home',
          tabBarIcon: ({ focused }) => <TabIcon icon="⌂" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="LeaderboardTab"
        component={LeaderboardScreen}
        options={{
          tabBarLabel: 'Leaderboard',
          tabBarIcon: ({ focused }) => <TabIcon icon="⊞" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="ProfileTab"
        component={ProfileScreen}
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
});
